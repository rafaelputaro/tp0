package com.putaro.app;

/**
 * Unidad a partir de la cual se puede construir una lista enlazada, de manera tal que contiene un 
 * objeto item y esta enlazado al siguiente nodo.
 */
public class Node 
{
	
	private Object item = null;
	
	private Node nextNode = null;
	
	/**
	 * @param item Objeto contenido en el nodo.
	 */
	public Node(Object item)
	{
		this.item = item;
	}
	
	public Object getItem()
	{
		return this.item;
	}
				
	public Node getNextNode()
	{
		return this.nextNode;
	}
	
	/**
	 *	Agrega un nodo hoja conteniendo el objeto parámetro.
	 * 
	 *	@param item Objeto que debe contener el nodo a agregar. 
	*/					
	public void addNextNode(Object item)
	{
		try
		{
			this.nextNode.addNextNode(item);
		} catch(Exception e)
		{
			this.nextNode = new Node(item);
		}
	}
}